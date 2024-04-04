package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

func prepareUUIDs(count int) (uuids []string) {

	for i := 0; i < count; i++ {
		uuids = append(uuids, uuid.New().String())
	}

	return
}

func RandomString(length int) string {
	const charset = `abcdefghijklmnopqrstuvwxyz`

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// prepare fake from addresses
func prepareFakeEmails(count int, domain string) (emails []string) {
	for i := 0; i < count; i++ {
		name := RandomString(7)     // Generate a random string of length 5 for the name.
		lastname := RandomString(8) // Generate a random string of length 5 for the lastname.
		id := rand.Intn(1000)       // Generate a random ID between 0 and 99.
		email := fmt.Sprintf("%s.%s.%d@%s", name, lastname, id, domain)
		emails = append(emails, email)
	}
	return
}

func prepareFakeSubjects(count int) (subjects []string) {
	for i := 0; i < count; i++ {
		subject := generateRandomText(3, 10)
		subjects = append(subjects, subject)
	}
	return
}

func generateRandomText(minSize int, maxSize int) string {
	words := []string{
		"urgent", "action", "required", "offer", "discount", "holiday", "sale", "alert", "notice",
		"update", "reminder", "exclusive", "limited", "opportunity", "chance", "win", "free",
		"premium", "access", "trial", "subscription", "upgrade", "enhanced", "improved", "new",
		"arrival", "launch", "introducing", "celebration", "thank", "you", "customer", "satisfaction",
		"feedback", "survey", "questionnaire", "invitation", "event", "webinar", "seminar", "conference",
		"workshop", "presentation", "newsletter", "announcement", "news", "update", "release", "version",
		"edition", "summary", "highlights", "features", "benefits", "advantages", "prospects", "solutions",
		"services", "products", "offerings", "selections", "choices", "varieties", "assortment", "portfolio",
		"catalog", "listings", "inventory", "stock", "availability", "accessibility", "convenience", "usability",
		"functionality", "performance", "efficiency", "effectiveness", "profitability", "savings", "earnings",
		"returns", "revenue", "sales", "deals", "bargains", "markdowns", "clearance", "closeouts", "liquidation",
		"fire", "sale", "auction", "bid", "offer", "proposal", "quote", "estimate", "inquiry", "interest",
		"curiosity", "exploration", "discovery", "learning", "education", "knowledge", "wisdom", "insight",
		"understanding", "awareness", "consciousness", "perception", "observation", "viewpoint", "perspective",
		"angle", "outlook", "vision", "foresight", "prediction", "forecast", "projection", "anticipation",
		"expectation", "speculation", "hypothesis", "theory", "conjecture", "guess", "estimate", "approximation",
		"calculation", "computation", "analysis", "evaluation", "assessment", "appraisal", "review", "examination",
		"investigation", "research", "study", "inquiry", "probe", "exploration", "survey", "scan", "check",
		"test", "trial", "experiment", "observation", "experience", "practice", "application", "exercise",
		"effort", "endeavor", "attempt", "try", "venture", "adventure", "feat", "achievement", "accomplishment",
		"success", "victory", "triumph", "conquest", "win", "gain", "advantage", "benefit", "reward", "prize",
		"award", "honor", "recognition", "acclaim", "praise", "commendation", "applause", "cheer", "ovation",
		"celebration", "party", "bash", "gathering", "assembly", "meeting", "rendezvous", "appointment", "date",
		"engagement", "commitment", "obligation", "responsibility", "duty", "task", "job", "work", "project",
		"assignment", "mission", "charge", "quest", "pursuit", "search", "hunt", "chase", "race", "competition",
		"contest", "challenge", "match", "game", "play", "sport", "recreation", "entertainment", "amusement",
		"diversion", "pastime", "hobby", "leisure", "relaxation", "rest", "break", "respite", "pause", "intermission",
		"interval", "gap", "space", "stretch", "span", "period", "duration", "term", "season", "time", "era",
		"age", "epoch", "generation", "cycle", "sequence", "series", "stream", "flow", "continuum", "progression",
		"development", "evolution", "growth", "expansion", "escalation", "increase", "rise", "surge", "wave",
		"trend", "direction", "movement", "motion", "dynamics", "force", "energy", "vigor", "vitality", "spirit",
		"enthusiasm", "zeal", "passion", "drive", "motivation", "ambition", "aspiration", "goal", "objective",
		"aim", "target", "purpose", "intent", "intention", "design", "plan", "strategy", "tactic", "method",
		"approach", "technique", "mode", "manner", "style", "fashion", "custom", "tradition", "convention",
		"practice", "routine", "habit", "tendency", "trend", "pattern", "norm", "standard", "criterion", "measure",
		"benchmark", "yardstick", "touchstone", "guideline", "rule", "regulation", "directive", "order", "decree",
		"command", "mandate", "edict", "proclamation", "announcement", "declaration", "statement", "pronouncement",
		"assertion", "affirmation", "testimony", "evidence", "proof", "confirmation", "verification", "validation",
		"authentication", "ratification", "endorsement", "approval", "acceptance", "consent", "agreement", "accord",
		"harmony", "concord", "unity", "unison", "synchronization", "coordination", "cooperation", "collaboration",
		"partnership", "alliance", "association", "relationship", "connection", "link", "bond", "tie", "attachment",
		"affiliation", "membership", "involvement", "participation", "engagement", "commitment", "dedication",
		"devotion", "loyalty", "fidelity", "allegiance", "fealty", "honor", "respect", "esteem", "admiration",
		"reverence", "veneration", "worship", "idolization", "adoration", "love", "affection", "attachment",
		"fondness", "liking", "preference", "taste", "inclination", "bent", "bias", "proclivity", "predisposition",
		"tendency", "penchant", "predilection", "partiality", "favor", "favoritism", "preference", "choice",
		"selection", "option", "alternative", "possibility", "opportunity", "chance", "prospect", "potential",
		"capability", "capacity", "ability", "competence", "skill", "talent", "aptitude", "faculty", "gift",
		"endowment", "asset", "resource", "advantage", "benefit", "value", "worth", "merit", "quality", "excellence",
		"superiority", "distinction", "uniqueness", "originality", "novelty", "innovation", "creativity", "imagination",
		"ingenuity", "inventiveness", "inspiration", "vision", "insight", "foresight", "perspicacity", "acumen",
		"wisdom", "knowledge", "understanding", "awareness", "consciousness", "perception", "cognition", "thought",
		"idea", "concept", "notion", "impression", "feeling", "sensation", "emotion", "mood", "temperament", "character",
		"personality", "identity", "individuality", "self", "ego", "essence", "substance", "nature", "essence",
		"core", "heart", "center", "nucleus", "hub", "focus", "epicenter", "axis", "pivot", "base", "foundation",
		"groundwork", "cornerstone", "keystone", "linchpin", "anchor", "mainstay", "backbone", "framework", "structure",
		"skeleton", "outline", "schema", "blueprint", "plan", "design", "pattern", "model", "prototype", "template",
		"mold", "form", "shape", "configuration", "composition", "constitution", "makeup", "build", "physique",
		"stature", "figure", "frame", "body", "anatomy", "physiology", "biology", "chemistry", "physics", "astronomy",
		"geology", "meteorology", "oceanography", "ecology", "environment", "nature", "universe", "cosmos", "world",
		"earth", "planet", "globe", "sphere", "realm", "domain", "field", "area", "territory", "region", "zone",
		"sector", "quarter", "district", "locality", "location", "place", "spot", "site", "setting", "scene",
		"landscape", "view", "panorama", "vista", "prospect", "outlook", "horizon", "background", "context", "framework",
		"matrix", "medium", "atmosphere", "ambience", "climate", "environment", "surroundings", "circumstances", "conditions",
		"situation", "status", "position", "stance", "standpoint", "perspective", "viewpoint", "angle", "approach",
		"method", "technique", "strategy", "tactic", "plan", "scheme", "plot", "design", "intent", "purpose", "aim",
		"objective", "goal", "target", "end", "ambition", "aspiration", "dream", "vision", "hope", "desire", "wish",
		"yearning", "longing", "craving", "hunger", "thirst", "appetite", "drive", "motivation", "impulse", "urge",
		"instinct", "inclination", "bent", "tendency", "propensity", "predisposition", "proclivity", "penchant",
		"predilection", "bias", "prejudice", "partiality", "favoritism", "preference", "liking", "taste", "appreciation",
		"affection", "attachment", "fondness", "love", "passion", "ardor", "enthusiasm", "zeal", "vigor", "energy",
		"vitality", "dynamism", "force", "power", "strength", "stamina", "endurance", "resilience", "fortitude",
		"courage", "bravery", "valor", "heroism", "prowess", "dexterity", "agility", "nimbleness", "flexibility",
		"mobility", "speed", "quickness", "rapidity", "haste", "dispatch", "efficiency", "effectiveness", "proficiency",
		"competence", "skill", "expertise", "mastery", "genius", "talent", "gift", "knack", "flair", "aptitude",
		"ability", "capacity", "capability", "potential", "possibility", "opportunity", "chance", "prospect", "future",
		"destiny", "fate", "lot", "fortune", "luck", "chance", "accident", "coincidence", "serendipity", "kismet",
		"destiny", "providence", "divine", "will", "fate", "predestination", "foreordination", "inevitability", "certainty",
		"necessity", "compulsion", "obligation", "duty", "responsibility", "requirement", "necessity", "need", "demand",
		"request", "plea", "appeal", "petition", "entreaty", "supplication", "prayer", "invocation", "benediction",
		"blessing", "grace", "mercy", "favor", "goodwill", "kindness", "generosity", "charity", "benevolence", "altruism",
		"philanthropy", "humanitarianism", "compassion", "sympathy", "empathy", "understanding", "tolerance", "patience",
		"forbearance", "leniency", "clemency", "mercy", "forgiveness", "pardon", "absolution", "exoneration", "acquittal",
		"release", "liberation", "freedom", "emancipation", "deliverance", "rescue", "salvation", "redemption", "recovery",
		"healing", "cure", "remedy", "solution", "answer", "resolution", "settlement", "reconciliation", "peace", "harmony",
		"accord", "agreement", "consensus", "unity", "cohesion", "solidarity", "cooperation", "collaboration", "partnership",
		"alliance", "association", "syndicate", "consortium", "coalition", "federation", "league", "union", "confederation",
		"confederacy", "empire", "kingdom", "realm", "domain", "territory", "province", "state", "nation", "country",
		"land", "region", "district", "county", "municipality", "city", "town", "village", "community", "society",
		"culture", "civilization", "heritage", "tradition", "custom", "practice", "ritual", "ceremony", "observance",
		"festival", "holiday", "celebration", "event", "occasion", "gathering", "assembly", "meeting", "conference",
		"convention", "symposium", "forum", "panel", "discussion", "debate", "dialogue", "conversation", "talk",
		"chat", "discourse", "monologue", "lecture", "speech", "address", "sermon", "homily", "oration", "proclamation",
		"announcement", "declaration", "statement", "testimony", "evidence", "proof", "documentation", "record",
		"report", "account", "narrative", "description", "explanation", "interpretation", "analysis", "assessment",
		"evaluation", "appraisal", "review", "critique", "criticism", "commentary", "opinion", "view", "perspective",
		"insight", "observation", "remark", "note", "annotation", "gloss", "comment", "reflection", "consideration",
		"thought", "idea", "concept", "notion", "perception", "impression", "feeling", "sensation", "emotion",
		"mood", "attitude", "stance", "position", "posture", "bearing", "demeanor", "manner", "behavior", "conduct",
		"action", "activity", "operation", "procedure", "process", "method", "technique", "practice", "routine",
		"habit", "custom", "tradition", "convention", "norm", "standard", "rule", "regulation", "guideline", "directive",
		"mandate", "order", "decree", "edict", "proclamation", "law", "statute", "ordinance", "regulation", "code",
		"precept", "principle", "doctrine", "creed", "belief", "conviction", "faith", "trust", "confidence", "assurance",
		"certainty", "security", "safety", "protection", "defense", "shield", "guard", "safeguard", "measure", "precaution",
		"preventive", "deterrent", "obstacle", "barrier", "hindrance", "impediment", "restriction", "limitation", "constraint",
		"control", "check", "curb", "brake", "restraint", "repression", "suppression", "inhibition", "discipline", "correction",
		"punishment", "penalty", "sanction", "fine", "charge", "fee", "cost", "price", "expense", "expenditure", "outlay",
		"investment", "spending", "budget", "allocation", "appropriation", "funding", "financing", "sponsorship", "support",
		"backing", "endorsement", "patronage", "assistance", "aid", "help", "relief", "rescue", "salvation", "deliverance",
		"liberation", "freedom", "emancipation", "independence", "autonomy", "sovereignty", "self-determination", "self-government",
		"democracy", "republic", "monarchy", "dictatorship", "tyranny", "despotism", "autocracy", "oligarchy", "theocracy",
		"anarchy", "chaos", "disorder", "confusion", "turmoil", "upheaval", "crisis", "emergency", "catastrophe", "disaster",
		"calamity", "tragedy", "misfortune", "adversity", "hardship", "difficulty", "challenge", "obstacle", "barrier",
		"hindrance", "impediment", "setback", "defeat", "failure", "loss", "deprivation", "lack", "shortage", "scarcity",
		"deficiency", "insufficiency", "want", "need", "necessity", "demand", "requirement", "obligation", "duty",
		"responsibility", "burden", "load", "weight", "pressure", "stress", "strain", "tension", "anxiety", "worry",
		"fear", "dread", "terror", "panic", "alarm", "shock", "horror", "apprehension", "suspense", "uncertainty",
		"doubt", "hesitation", "reluctance", "indecision", "vacillation", "wavering", "fluctuation", "variation", "change",
		"alteration", "modification", "adjustment", "adaptation", "transformation", "conversion", "revolution", "reform",
		"innovation", "renovation", "restoration", "renewal", "revival", "resurgence", "rebirth", "renaissance", "reawakening",
		"regeneration", "rejuvenation", "refreshment", "replenishment", "recovery", "recuperation", "healing", "cure",
		"remedy", "solution", "resolution", "settlement", "reconciliation", "agreement", "compromise", "bargain", "deal",
		"contract", "pact", "treaty", "accord", "concord", "harmony", "peace", "truce", "ceasefire", "armistice", "detente",
		"rapprochement", "cooperation", "collaboration", "partnership", "alliance", "association", "syndicate", "consortium",
		"coalition", "federation", "league", "union", "confederation", "confederacy", "conglomerate", "corporation", "company",
		"firm", "enterprise", "business", "organization", "institution", "establishment", "agency", "bureau", "office",
		"department", "division", "unit", "section", "branch", "wing", "arm", "squad", "team", "group", "crew", "band",
		"gang", "mob", "pack", "herd", "flock", "swarm", "school", "colony", "community", "society", "association", "club",
		"circle", "league", "guild", "union", "federation", "confederation", "alliance", "coalition", "bloc", "axis",
		"partnership", "consortium", "syndicate", "network", "system", "complex", "aggregate", "whole", "ensemble", "combination",
		"mixture", "blend", "fusion", "integration", "unification", "amalgamation", "merger", "union", "consolidation",
		"synthesis", "composition", "construction", "fabrication", "manufacture", "production", "creation", "invention",
		"innovation", "development", "growth", "expansion", "escalation", "increase", "rise", "surge", "boom", "burst",
		"explosion", "outbreak", "eruption", "flare-up", "upswing", "upturn", "upheaval", "revolution", "transformation",
		"transition", "changeover", "shift", "switch", "turnaround", "reversal", "overhaul", "makeover", "remodeling",
		"restructuring", "reorganization", "reform", "revision", "correction", "adjustment", "modification", "alteration",
		"amendment", "update", "upgrade", "improvement", "enhancement", "refinement", "polishing", "perfecting", "optimization",
		"streamlining", "simplification", "clarification", "elucidation", "explanation", "interpretation", "translation",
		"conversion", "transcription", "transliteration", "decoding", "deciphering", "unraveling", "solving", "resolving",
		"answering", "clarifying", "enlightening", "illuminating", "revealing", "unveiling", "disclosing", "exposing",
		"uncovering", "unmasking", "demonstrating", "proving", "confirming", "verifying", "validating", "authenticating",
		"ratifying", "endorsing", "approving", "sanctioning", "authorizing", "licensing", "permitting", "allowing",
		"enabling", "facilitating", "empowering", "equipping", "arming", "strengthening", "fortifying", "bolstering",
		"supporting", "backing", "upholding", "sustaining", "maintaining", "preserving", "protecting", "defending",
		"shielding", "guarding", "safeguarding", "securing", "ensuring", "guaranteeing", "warranting", "pledging",
		"vowing", "promising", "committing", "dedicating", "devoting", "giving", "donating", "contributing", "offering",
		"presenting", "awarding", "bestowing", "granting", "conferring", "imparting", "transferring", "conveying",
		"delivering", "handing", "passing", "dispatching", "sending", "mailing", "posting", "shipping", "transporting",
		"conveying", "moving", "transferring", "shifting", "relocating", "transplanting", "migrating", "traveling",
		"journeying", "voyaging", "touring", "exploring", "venturing", "adventuring", "roaming", "wandering", "drifting",
		"meandering", "rambling", "strolling", "walking", "hiking", "trekking", "climbing", "scaling", "ascending",
		"mountaineering", "alpinism", "skiing", "snowboarding", "surfing", "swimming", "diving", "snorkeling", "fishing",
		"hunting", "shooting", "archery", "fencing", "wrestling", "boxing", "martial", "arts", "karate", "judo",
		"taekwondo", "kung", "fu", "aikido", "jiu-jitsu", "sumo", "wushu", "capoeira", "kickboxing", "muay", "thai"}

	// Determine the number of words in the subject (between minSize and maxSize)
	numWords := rand.Intn(maxSize-minSize) + minSize

	var subjectWords []string
	for i := 0; i < numWords; i++ {
		// Select a random word and add it to the subjectWords slice
		word := words[rand.Intn(len(words))]
		subjectWords = append(subjectWords, word)
	}

	// Combine the selected words into a single string
	text := strings.Join(subjectWords, " ")

	return text
}
